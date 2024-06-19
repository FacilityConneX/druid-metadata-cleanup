package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func main() {
	//flags
	configFileFlag := flag.String("config-file", "config.yaml", "Configuration file for tool.")
	showTasksFlag := flag.Bool("show-tasks", false, "Include all tasks found for the period in the preview.")
	deleteFlag := flag.Bool("delete", false, "Delete items from S3 and SQL after preview.")
	flag.Parse()

	viper.SetConfigFile(*configFileFlag)
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	objectKeys, err := preview(*showTasksFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		panic("error with preview")
	}

	if *deleteFlag {
		err := delete(objectKeys)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			panic("error with preview")
		}
	}

	db.Close()
}

// preview will always run
func preview(show bool) ([]string, error) {
	dbURL := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v",
		viper.Get("database.user"),
		viper.Get("database.password"),
		viper.Get("database.host"),
		viper.Get("database.port"),
		viper.Get("database.db"),
	)

	err := initDB(dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil, err
	}

	endTime := viper.GetString("endTime")
	tasks, err := queryMetadata(endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil, err
	}

	objectKeys := []string{}
	for _, v := range tasks {
		if show {
			fmt.Println(v.id, "  -  ", v.createdDate)
		}
		keyBase := viper.Get("s3.keyBase")
		logKey := fmt.Sprintf("%v/%v/log", keyBase, v.id)
		reportKey := fmt.Sprintf("%v/%v/report.json", keyBase, v.id)
		objectKeys = append(objectKeys, logKey, reportKey)
	}

	fmt.Println(len(tasks), "to delete from SQL.", len(objectKeys), "to delete from S3.")
	return objectKeys, nil
}

// Only run when `--delete` flag is passed
func delete(objects []string) error {
	initS3Client()

	endTime := viper.GetString("endTime")

	bucket := viper.GetString("s3.bucket")
	batchSize := viper.GetInt("s3.batchSize")

	fmt.Println("Deleting from S3.")
	err := batchDeleteObjects(bucket, objects, batchSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "S3 delete error: %v\n", err)
		return err
	}

	fmt.Println("Deleting from SQL.")
	err = deleteMetadata(endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SQL delete error: %v\n", err)
		return err
	}
	return nil
}