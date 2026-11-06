CREATE TABLE role_accesses (
       role_id INT,
       access_id INT,
       PRIMARY KEY (role_id, access_id),
       FOREIGN KEY (role_id) REFERENCES roles(role_id),
       FOREIGN KEY (access_id) REFERENCES accesses(access_id)
);
