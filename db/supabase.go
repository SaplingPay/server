package db

import (
	"log"

	storage_go "github.com/supabase-community/storage-go"
)

var Supabase *storage_go.Client

func ConnectSupabase(supabaseUrl string, supabaseKey string) {
	// Set client options
	// supabase := supa.CreateClient(supabaseUrl, supabaseKey)
	supabase := storage_go.NewClient(supabaseUrl+".supabase.co/storage/v1", supabaseKey, nil)

	Supabase = supabase

	log.Println("Connected to Supabase!")
}
