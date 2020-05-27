package copier_test

import (
	"encoding/json"
	"testing"

	"taurus-admin/library/utils/copier"
)

func BenchmarkCopyStruct(b *testing.B) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		copier.Copy(&Employee{}, &user)
	}
}

func BenchmarkNameCopy(b *testing.B) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		employee := &Employee{
			Name:      Name,
			Nickname:  &Nickname,
			Age:       int64(Age),
			FakeAge:   int(*FakeAge),
			DoubleAge: DoubleAge(),
			Notes:     Notes,
		}
		Role(Role)
	}
}

func BenchmarkJsonMarshalCopy(b *testing.B) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		data, _ := json.Marshal(user)
		var employee Employee
		json.Unmarshal(data, &employee)

		DoubleAge = DoubleAge()
		Role(Role)
	}
}
