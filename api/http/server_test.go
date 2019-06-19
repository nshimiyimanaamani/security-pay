package http

import "testing"

func TestRunServer(t *testing.T){
	//create a new server instace
	srv:= New()

	if err:=srv.Run(); err!=nil{
		t.Errorf("unexpected error '%v'", err)
	}
}