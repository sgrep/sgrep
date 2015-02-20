package sgrep

import "io/ioutil"
import "path"
import "fmt"

type Directory struct {
    // The name of this directory relative to its parent directory.
    // Eg., for structure a/b/c/, would have a directory with name
    // a, name b, and name c.
    Name string
    Files []string
    Directories []* Directory
}

func (dir* Directory) PrettyPrint() {
	dir.pretty_print_helper(0)
}


/**
  Returns all files as fully-qualified filename from directory dir.
*/
func (dir* Directory) ListFiles() []string {
    to_return := make([]string,0)

    for _, filename := range dir.Files {
        fq_filename := path.Join(dir.Name,filename)
        to_return = append(to_return, fq_filename)
    }

    for _, subdir := range dir.Directories {
        subdir_file_slice := subdir.ListFiles()

        for _, filename := range subdir_file_slice {
            fq_filename := path.Join(dir.Name,filename)
            to_return = append(to_return, fq_filename)
        }
    }
    return to_return
}

func (dir *Directory) pretty_print_helper(indentation_level uint32) {

    indent_str := ""
    for i := uint32(0); i < indentation_level; i++ {
        indent_str += "\t"
    }
    fmt.Println(indent_str + dir.Name + "/")
    
    for _, filename := range dir.Files {
        fmt.Println(indent_str + "\t" + filename)
    }
        
    for _, sub_directory := range dir.Directories {
        sub_directory.pretty_print_helper(indentation_level + 1)
    }
}


/**
  @param dir_to_walk The directory that we want to list all
  subdirectories of.

  @returns A Directory struct containing all subfiles and folders.
*/
func WalkFolder( dir_to_walk string) * Directory {
    dir_contents_list, err := ioutil.ReadDir(dir_to_walk)
    if err != nil {
        panic ("IOError when reading directory " + dir_to_walk)
    }

	root := new(Directory)
    root.Name = path.Base(dir_to_walk)
    
	for _, file_info := range dir_contents_list {
        fully_qualified_path := path.Join(dir_to_walk,file_info.Name())
        if err != nil {
            panic(
                "Could not stat file or folder named " +
                    fully_qualified_path)
        }
        
        if file_info.IsDir() {
            sub_directory := WalkFolder(fully_qualified_path)
            root.Directories = append(root.Directories,sub_directory)
        } else {
            root.Files = append(root.Files, file_info.Name())
        }
    }
    return root
}
