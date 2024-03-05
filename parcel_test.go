package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел
	// Для повышения уникальности в качестве seed используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// Создаём тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	//добавляю новую посылку в БД, проверка на отсутствие ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEqual(t, 0, id)

	// Получаю только что добавленную посылки
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	// проверка, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	require.Equal(t, parcel.Address, storedParcel.Address)
	require.Equal(t, parcel.Client, storedParcel.Client)
	require.Equal(t, parcel.CreatedAt, storedParcel.CreatedAt)
	require.Equal(t, parcel.Status, storedParcel.Status)

	// удаляю добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(id)
	require.NoError(t, err)

	// проверяю, что посылку больше нельзя получить из БД
	_, err = store.Get(id)
	require.Error(t, err) // Ожидаю ошибку при попытке получить посылку
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	
	// Добавляю новую посылку в БД, проверяю отсутствие ошибки и наличии идентификатора
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEqual(t, 0, id)
	
	// обновляю адрес
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// получаю добавленную посылку и проверяю, что адрес обновился
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, storedParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	//Добавляю новую посылку в БД, проверяю отсутствие ошибки и наличии идентификатора
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEqual(t, 0, id)

	//Обновляю статус
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	// получаю добавленную посылку и проверяю, что статус обновился
	parcelStore, err := store.Get(id)
	require.NoError(t, err)
	require.NotEqual(t, parcelStore.Status, ParcelStatusRegistered)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcelsSlice := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	//Задаю всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)

	for i := 0; i < len(parcelsSlice); i++ {
		parcelsSlice[i].Client = client
		id, err := store.Add(parcelsSlice[i])
		require.NoError(t, err)
		require.NotEqual(t, 0, id)

		// обновляем идентификатор добавленной у посылки
		parcelsSlice[i].Number = id

		// сохраняем добавленную посылку в структуру map по ID, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcelsSlice[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	// проверка, что количество полученных посылок совпадает с количеством добавленных
	require.Equal(t, len(parcelsSlice), len(storedParcels))

	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// проверка, что все посылки из storedParcels есть в parcelMap
		value, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		// проверка, что значения полей полученных посылок заполнены верно
		require.Equal(t, value, parcel)
	}
}
