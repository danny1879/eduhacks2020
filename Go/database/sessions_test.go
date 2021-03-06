package database

import (
	"testing"
)

type DecryptedTest struct {
	Count       int
	Login       bool
	CipherText  string
	SessionName string
}

const (
	cipherText1 = "MTU5OTI4NTcyNHxEdi1CQkFFQ180SUFBUkFCRUFBQU5mLUNBQUlHYzNSeWFXNW5EQWNBQld4dloybHVCR0p2YjJ3Q0FnQUJCbk4wY21sdVp3d0hBQVZqYjNWdWRBTnBiblFFQWdBQXxmCAprCZYCD-vgCRVPdvSPmFoH_nPW8fGA3RN_kumtIQ=="
	cipherText2 = "MTU5OTI4NjMyNHxEdi1CQkFFQ180SUFBUkFCRUFBQU5mLUNBQUlHYzNSeWFXNW5EQWNBQld4dloybHVCR0p2YjJ3Q0FnQUJCbk4wY21sdVp3d0hBQVZqYjNWdWRBTnBiblFFQWdBR3wZKOsOY-U6j_L60ikqU2mk9vF9tSH5taFm_FskYpn_jw=="
)

var DecryptedTests = []DecryptedTest{
	{0, true, cipherText1, SessionName},
	{3, true, cipherText2, SessionName},
}

//func TestSessionManager_DecryptedData(t *testing.T) {
//	for _, test := range DecryptedTests {
//		session := SessionManager{}
//		session.DecryptedData(test.CipherText, test.SessionName)
//		if session.Values["count"] != test.Count {
//			t.Errorf("Count = %v, want %v", session.Values["count"], test.Count)
//		}
//		if session.Values["login"] != test.Login {
//			t.Errorf("Login = %v, want %v", session.Values["login"], test.Login)
//		}
//	}
//}

func TestSessionManager_EncryptedData(t *testing.T) {
	session := SessionManager{Values: make(map[interface{}]interface{})}
	session.Values["count"] = 300
	session.Values["login"] = true

	cipherText, err := session.EncryptedData(SessionName)
	if err != nil {
		t.Errorf("Encrypted an error has occurred %s", err.Error())
	}
	session.Values["count"] = 10
	session.Values["login"] = false
	session.DecryptedData(cipherText, SessionName)
	if session.Values["count"] != 300 {
		t.Errorf("Count = %v, want %v", session.Values["count"], 300)
	}
	if session.Values["login"] != true {
		t.Errorf("Login = %v, want %v", session.Values["count"], true)
	}
}

//func TestSessionManager_SaveData(t *testing.T) {
//	session := SessionManager{Values: make(map[interface{}]interface{})}
//	session.Values["count"] = 8000
//	session.Values["login"] = true
//
//	cipherText, err := session.EncryptedData(SessionName)
//	if err != nil {
//		t.Errorf("Encrypted an error has occurred %s", err.Error())
//	}
//	session.SaveData("5f5321edbc36da097814b7da", cipherText)
//}

//func TestSessionManager_GetData(t *testing.T) {
//	setting.DefaultSetting()
//	type fields struct {
//		Values map[interface{}]interface{}
//	}
//	type args struct {
//		id string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    interface{}
//		wantErr bool
//	}{
//		{"test", fields{Values: make(map[interface{}]interface{})}, args{id: "5f6853c5f742920128fd123d"}, nil, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			se := &SessionManager{
//				Values: tt.fields.Values,
//			}
//			got, err := se.GetData(tt.args.id)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetData() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
