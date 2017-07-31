database "alice_db" {
  endpoint = "alice.example.com"
  port = 5432
  user = "admin"
  password = "admin"
  database_name = "apple"
}

group "dev" {
  policies = [
    "alice_db_readonly",
  ]
}

policy "alice_db_readonly" {
  database = "alice_db"

  queries = [
    "CREATE ROLE {{ .Name }} WITH LOGIN ENCRYPTED PASSWORD {{ .Password }};",
    "GRANT SELECT ON ALL TABLES IN SCHEMA public TO {{ .Name }};",
    "GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO {{ .Name }};"
  ]
}

user "alice" {
  group = "dev"
}
