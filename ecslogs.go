// App: ecslogs
// Author: Will Houle
// Date: May 2017
// Description: log collection app for ECS Container Agent; to replace ecs-log-collector shell script

package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// Gets current directory
	// Golang 1.8 recommends using os.Executable
	currdir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// sets target directory for outputting gathered information "infodir"
	infodir := currdir + "/collect"
	infosystem := infodir + "/system" // TODO: remove this dir, think it's unnecessary

	fmt.Println("current directory:", currdir)
	fmt.Println("target directory:", infodir)
	fmt.Println("System dir:", infosystem)
	CreateArchive()
	getEc2InstanceID()
}

// CreateArchive - packs up the gathered log files into a .tgz archive
func CreateArchive() {
	// output
	file, err := os.Create("collect.tar.gz")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Set up writer
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// paths to be included into the archive
	paths := []string{
		"test.log",
	}

	// add files to archive
	for i := range paths {
		if err := addFile(tw, paths[i]); err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("CreateArchive function end")
}

func addFile(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if stat, err := file.Stat(); err == nil {
		// now lets create the header as needed for this file within the tarball
		header := new(tar.Header)
		header.Name = path
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		// write the header to the tarball archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// copy the file data to the tarball
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

// GetEc2InstanceId - uses AWS SDK to get the instance ID
func getEc2InstanceID() string {
	//svc := ec2metadata.New()
	//fmt.Println(svc)
	instanceID := "i-a1b2c3d4e5f6g7"
	fmt.Println("Instance ID:", instanceID)
	return instanceID
}
