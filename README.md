# **GoVCS - A Simple Version Control System in Go**
GoVCS is a lightweight, Git-like version control system built using **Go**. It uses **Merkle trees** to efficiently track file changes and **SQLite** as a database backend for storing commit history.  

This project is a foundational step toward building a more advanced VCS, with future enhancements such as **multi-user collaboration**, **branching**, and **merge conflict resolution**.  

---

## **📌 How Merkle Tree Works in GoVCS**
In GoVCS, the **Merkle tree** plays a crucial role in securing and tracking file versions efficiently.  

1. **File Hashing**  
   - Each file's content is hashed using SHA-256.  
2. **Constructing the Merkle Tree**  
   - The hashes of all committed files are combined to form a tree.  
   - Leaf nodes contain individual file hashes, while parent nodes store the hash of concatenated child nodes.  
3. **Root Hash as a Unique Commit ID**  
   - The final **Merkle Root Hash** represents the entire commit.  
   - If any file content changes, the root hash will also change, ensuring data integrity.  

📌 **Example Merkle Tree for 4 files:**  

```
         Root Hash
       /          \
  Hash1_2       Hash3_4
   /    \         /    \
H(file1) H(file2) H(file3) H(file4)
```
- If any file is modified, the root hash will change, making it easy to detect differences between commits.

---

## **📂 Project Structure**
```
go-vcs/
│── merkle/               # Merkle Tree logic
│   ├── tree.go
│── Repository/               # Repo files exists here
│   ├── file1.txt
    ├── file2.txt
    .
    .
    .
│── vcs/                  # Version Control System logic
│   ├── storage.go        # Database interaction
│   ├── vcs.go            # Core VCS commands
│── main.go               # CLI entry point
│── go.mod                # Go module dependencies
│── go.sum                # Dependency checksums
│── README.md             # Project documentation
```

---

## **🚀 How to Run the Project**
### **📌 Prerequisites**
- Install **Go**: [Download Go](https://go.dev/dl/)
- Install **SQLite3**
- Clone this repository:
  ```sh
  git clone https://github.com/your-username/go-vcs.git
  cd go-vcs
  ```

### **📌 Build the project**
```sh
go build -o go-vcs
```
This will create an executable file named `go-vcs` (or `go-vcs.exe` on Windows).

### **📌 Initialize the repository**
```sh
./go-vcs init
```
Creates an empty **vcs.db** file.

---

## **🔨 Commands to Test GoVCS**
### **1️⃣ Add and Commit Files**
```sh
./go-vcs commit "custom message"
```
- Reads and stores file contents in the database.
- Creates a Merkle tree to generate a unique commit ID.

### **2️⃣ View Commit History**
```sh
./go-vcs history
```
- Displays all past commits.

### **3️⃣ Revert to an Old Commit**
```sh
./go-vcs revert <commit_id>
```
- Restores all files to the state in the specified commit.

---

## **🛠️ Future Enhancements**
### **🔹 Multi-User Collaboration**
- Implement **user authentication** so multiple users can work on the same repository.
- Track which user made a commit.

### **🔹 Branching and Merging**
- Enable users to create and switch between branches.
- Implement a **merge strategy** to combine changes from different branches.

### **🔹 Merge Conflict Resolution**
- When two users modify the same file, detect conflicts and allow **manual resolution**.

### **🔹 Remote Repositories (Like GitHub)**
- Implement a remote server for pushing/pulling changes.
- Synchronize commits across different machines.

---

## **💡 Conclusion**
GoVCS is a simple but powerful version control system that lays the foundation for more advanced features. With Merkle trees ensuring data integrity and SQLite handling commit history, it provides a great starting point for understanding how Git-like systems work.  

🚀 **Want to contribute?** Fork the repo and submit a PR!  
📢 **Feedback?** Open an issue or reach out!  

