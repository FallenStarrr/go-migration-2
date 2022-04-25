package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func MakeFolder() {
	os.Mkdir("C:/Users/u14911/Desktop/Maquette/Миграции/" + os.Args[1], os.ModePerm);
}

	/*
			1 arg: folder & file name & name
			2 arg: parent_id
			3 arg: deepest_node
    */
func MakeJSON() {
	my_json,err := os.Create("C:/Users/u14911/Desktop/Maquette/Миграции/"+os.Args[1] + "/" + os.Args[1] + ".json")
	check(err)

	json_content := fmt.Sprintf(
		`[
		 {
		 "name": "%v",
		 "parent_id": "%v",
		 "deepest_node": "%v"
		 }
		 ]
		`, 
	os.Args[1], 
	os.Args[2], 
	os.Args[3])

	my_json.WriteString(json_content)

	fmt.Println(json_content)
}



func main() {
	MakeFolder()
	MakeJSON()

	sql_migration, err1 := os.Create("C:/Users/u14911/Desktop/Maquette/Миграции/" + os.Args[1] + "/" + os.Args[1] + ".sql")
	meta, err := os.Create("C:/Users/u14911/Desktop/Maquette/Миграции/" + os.Args[1] + "/" + os.Args[1] + ".txt")
	check(err)
	check(err1)
	type docs []struct {
		Name        string `json:"name"`
		ParentID    string `json:"parent_id"`
		DeepestNode bool   `json:"deepest_node"`
	}
	
	filea,_ := ioutil.ReadFile("C:/Users/u14911/Desktop/Maquette/Миграции/" + os.Args[1] + "/" + os.Args[1] + ".json")

	ads:= docs{}
	json.Unmarshal(filea,&ads)


	defer sql_migration.Close()

	query := ""
	for _, v := range ads {
		if v.ParentID == "" {
			v.ParentID = "null"
		}else {
			v.ParentID = "'"+v.ParentID+"'"
		}
		id := uuid.NewString()
		query += fmt.Sprintf(
			`

				 insert into document_type (id, doc_type, parent_id, deepest_node, expiration_time)
				 values (%q, %q, %v, %v, 99);
                 insert into roles(id, role)
				 values (%q,'COMMON-READ');
				 insert into roles(id, role)
				 values (%q,'COMMON-WRITE');

				

				 `,   id,  v.Name, v.ParentID, v.DeepestNode,  id,id)
				 text := fmt.Sprintf("Doc_Type: %v\n -------------  \n Name:  %v\n -------------   \nParent:  %v\n ------------- \n  Deepest Node:  %v\n", id,   v.Name,  v.ParentID,  v.DeepestNode) 
				 meta.WriteString(text)
	}
	query = strings.Replace(query, "\"", "'", -1)
	fmt.Println(query)
	sql_migration.WriteString(query)

	
}
