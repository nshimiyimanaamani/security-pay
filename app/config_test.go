package app

import "testing"
func TestConfig(t *testing.T){
	config, err:= NewConfig()
	
	t.Run("must return nil error", func(t *testing.T){
		if err!= nil{
			t.Errorf("unexpected error:'%v'", err)
		}
	})

	t.Run("config contains a private key", func(t *testing.T){
		if config.GetPrivateKey() ==nil{
			t.Errorf("private key must be set")
		}
	})

	t.Run("config contains a public key", func(t *testing.T){
		if config.GetPublicKey() ==nil{
			t.Errorf("public key must be set")
		}
	})
}