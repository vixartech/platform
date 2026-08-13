insert into users (
  username,
  password_login,
  password,
  email,
  name,
  type,
  user_role_id,
  list_role_id,
  status
) values (
  :'api_user',
  false,
  :'token_hash',
  :'api_user' || '@api',
  :'api_user',
  'api',
  1,
  null,
  'enabled'
)
on conflict (username) do update set
  password_login = excluded.password_login,
  password = excluded.password,
  email = excluded.email,
  name = excluded.name,
  type = excluded.type,
  user_role_id = excluded.user_role_id,
  list_role_id = excluded.list_role_id,
  status = excluded.status,
  updated_at = now();

insert into settings (key, value)
values (
  'smtp',
  '[{
    "name": "mailpit",
    "uuid": "",
    "enabled": true,
    "host": "platform-mailpit",
    "hello_hostname": "",
    "port": 1025,
    "auth_protocol": "none",
    "username": "",
    "password": "",
    "email_headers": [],
    "max_conns": 10,
    "max_msg_retries": 1,
    "msg_retry_delay": "10ms",
    "idle_timeout": "15s",
    "wait_timeout": "5s",
    "tls_type": "none",
    "tls_skip_verify": false,
    "from_addresses": ["Vixar <no-reply@vixar.cc>"]
  }]'::jsonb
)
on conflict (key) do update set
  value = excluded.value,
  updated_at = now();
