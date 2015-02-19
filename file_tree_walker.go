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

func (dir *Directory) pretty_print_helper(indentation_level uint32) {

    indent_str := ""
    for i := uint32(0); i < indentation_level; i++ {
        indent_str += "\t"
    }
    fmt.Println(indent_str + dir.Name)
    
    for _, filename := range dir.Files {
        fmt.Println(indent_str + filename)
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
func walk_folder( dir_to_walk string) * Directory {
    dir_contents_list, err := ioutil.ReadDir(dir_to_walk)
    if err != nil {
        panic ("IOError when reading directory " + dir_to_walk)
    }

	root := new(Directory)
    root.Name, _ = path.Split(dir_to_walk)
    
	for _, file_info := range dir_contents_list {
        fully_qualified_path := path.Join(dir_to_walk,file_info.Name())
        if err != nil {
            panic(
                "Could not stat file or folder named " +
                    fully_qualified_path)
        }
        
        if file_info.IsDir() {
            sub_directory := walk_folder(fully_qualified_path)
            root.Directories = append(root.Directories,sub_directory)
        } else {
            root.Files = append(root.Files, file_info.Name())
        }
    }
    return root
}
