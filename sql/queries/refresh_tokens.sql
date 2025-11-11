-- name: CreateRefreshToken :one
insert into
    refresh_tokens (
        token,
        created_at,
        updated_at,
        user_id,
        expires_at
    )
values
    ($1, NOW(), NOW(), $2, $3) returning *;

-- name: RevokeRefreshToken :one
update
    refresh_tokens
set
    updated_at = now(),
    revoked_at = now()
where
    token = $1 returning *;

-- name: GetUserFromRefreshToken :one

select users.* from users join refresh_tokens on users.id = refresh_tokens.user_id
where refresh_tokens.token = $1
  and refresh_tokens.revoked_at is null
  and refresh_tokens.expires_at > now();