policy "alice_db_writable" {
  database = "alice_db"

  queries = [
    "CREATE ROLE {{ .Name }} WITH LOGIN ENCRYPTED PASSWORD {{ .Password }};",
    "GRANT ALL ON ALL TABLES IN SCHEMA public TO {{ .Name }};",
    "GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO {{ .Name }};"
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

policy "bob_db_writable" {
  database = "bob_db"

  queries = [
    "CREATE ROLE {{ .Name }} WITH LOGIN ENCRYPTED PASSWORD {{ .Password }};",
    "GRANT ALL ON ALL TABLES IN SCHEMA public TO {{ .Name }};",
    "GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO {{ .Name }};"
  ]
}

policy "bob_db_readonly" {
  database = "bob_db"

  queries = [
    "CREATE ROLE {{ .Name }} WITH LOGIN ENCRYPTED PASSWORD {{ .Password }};",
    "GRANT SELECT ON ALL TABLES IN SCHEMA public TO {{ .Name }};",
    "GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO {{ .Name }};"
  ]
}

policy "carol_db_writable" {
  database = "carol_db"

  queries = [
    "CREATE ROLE {{ .Name }} WITH LOGIN ENCRYPTED PASSWORD {{ .Password }};",
    "GRANT ALL ON ALL TABLES IN SCHEMA public TO {{ .Name }};",
    "GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO {{ .Name }};"
  ]
}
