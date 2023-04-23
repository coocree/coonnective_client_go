package api

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"sync"
	"time"
)

// StorageModel é responsável por fornecer uma interface para armazenar e recuperar dados
// usando o BoltDB e criptografando as informações armazenadas.
type StorageModel struct {
	name     string
	password string
	db       *bolt.DB
	security SecurityModel
}

var instanceStorage *StorageModel
var once sync.Once

// Storage retorna a instância única de StorageModel.
func Storage() *StorageModel {
	once.Do(func() {
		_name := "Storage"
		_password := "api-storage"
		db, err := bolt.Open(getDatabaseFilename(_name), 0600, nil)

		if err != nil {
			apiErrorModel := ApiError([]string{"Falha ao localizar local storage"}, "Storage", "Storage",
				"LOCAL_STORAGE_FAILURE", time.Now(), []string{_name, _password})
			E(apiErrorModel.ToString(), apiErrorModel.Code, nil)
			panic(apiErrorModel.Code)
		}

		instanceStorage = &StorageModel{
			name:     _name,
			password: _password,
			db:       db,
			security: Security(),
		}
	})

	return instanceStorage
}

// Add adiciona um valor no armazenamento.
// Recebe uma chave (key) e um valor (value) como argumentos.
// O valor é criptografado antes de ser armazenado.
func (s *StorageModel) Add(key string, value string) {
	var apiErrorModel ErrorModel
	if len(key) == 0 || len(value) == 0 {
		apiErrorModel = ApiError([]string{"Chave ou valor vazio"}, "Storage", "Add",
			"STORAGE_EMPTY_KEY_OR_VALUE", time.Now(), value)
	}

	errUpdate := s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.name))
		if err != nil {
			return err
		}

		encryptedKey := s.security.EncodeMD5(key + s.password)
		encryptedValue, _ := s.security.Encrypt(value, s.password)
		err = bucket.Put([]byte(encryptedKey), []byte(encryptedValue))
		return err
	})
	if errUpdate != nil {
		apiErrorModel = ApiError([]string{"Falha ao adicionar valor no storage", errUpdate.Error()}, "Storage",
			"Add", "STORAGE_ADD_ERROR", time.Now(), value)
	}

	if apiErrorModel.HasError() {
		E(apiErrorModel.ToString(), apiErrorModel.Code, nil)
		panic(apiErrorModel.Code)
	}
}

// Remove remove um valor do armazenamento.
// Recebe uma chave (key) como argumento e remove o valor correspondente do armazenamento.
func (s *StorageModel) Remove(key string) {
	var apiErrorModel ErrorModel
	if len(key) == 0 {
		apiErrorModel = ApiError([]string{"Chave vazio"}, "Storage", "Remove",
			"STORAGE_EMPTY_KEY", time.Now(), key)
	}

	errUpdate := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.name))
		if bucket == nil {
			return fmt.Errorf("Falha ao adicionar valor no storage")
		}

		encryptedKey := s.security.EncodeMD5(key + s.password)
		err := bucket.Delete([]byte(encryptedKey))
		return err
	})
	if errUpdate != nil {
		apiErrorModel = ApiError([]string{errUpdate.Error()}, "Storage",
			"Remove", "STORAGE_REMOVE_ERROR", time.Now(), key)
	}

	if apiErrorModel.HasError() {
		E(apiErrorModel.ToString(), apiErrorModel.Code, nil)
		panic(apiErrorModel.Code)
	}
}

// Read lê um valor do armazenamento.
// Recebe uma chave (key) como argumento e retorna o valor descriptografado correspondente.
// Retorna um erro se a chave não for encontrada.
func (s *StorageModel) Read(key string) string {
	var apiErrorModel ErrorModel
	if len(key) == 0 {
		apiErrorModel = ApiError([]string{"Chave vazio"}, "Storage", "Read",
			"STORAGE_EMPTY_KEY", time.Now(), key)
	}

	var result string
	errUpdate := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.name))
		if bucket == nil {
			return fmt.Errorf("bucket não encontrado")
		}
		encryptedKey := s.security.EncodeMD5(key + s.password)
		encryptedValue := bucket.Get([]byte(encryptedKey))
		if encryptedValue == nil {
			return fmt.Errorf("chave não encontrada")
		}
		result, _ = s.security.Decrypt(string(encryptedValue), s.password)
		return nil
	})

	if errUpdate != nil {
		apiErrorModel = ApiError([]string{errUpdate.Error()}, "Storage",
			"Read", "STORAGE_READ_ERROR", time.Now(), key)
	}

	if apiErrorModel.HasError() {
		E(apiErrorModel.ToString(), apiErrorModel.Code, nil)
		panic(apiErrorModel.Code)
	}

	return result
}

// Close libera recursos e define o armazenamento como nulo.
func (s *StorageModel) Close() error {
	return s.db.Close()
}

// getDatabaseFilename retorna o nome do arquivo de banco de dados para o nome fornecido.
func getDatabaseFilename(name string) string {
	hasher := sha1.New()
	hasher.Write([]byte(name))
	return hex.EncodeToString(hasher.Sum(nil)) + ".db"
}
