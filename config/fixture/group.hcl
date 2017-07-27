group "dev" {
  policies = [
    "alice_db_readonly",
    "bob_db_readonly",
    "carol_db_writable"
  ]
}

group "core" {
  policies = [
    "alice_db_writable",
    "bob_db_writable",
    "carol_db_writable"
  ]
}
