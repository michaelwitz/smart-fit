# Database Migrations

## Current Status: Empty (Development Phase)

This directory is intentionally empty during development. We follow a **no-migrations philosophy** where we:

1. **Recreate the database** on schema changes
2. **Use complete schema files** in `../schemas/complete_database_schema.sql`
3. **Reseed with fresh data** from `../seeds/`

## When Migrations Will Be Added

Migrations will be introduced when the app moves to production and we need to:

- Preserve existing production data
- Track schema evolution over time
- Enable rollback capabilities
- Support multiple environments with different schema versions

## Current Schema Management

- **Source of Truth**: `../schemas/complete_database_schema.sql`
- **Setup Process**: `../../scripts/setup-database.sh`
- **Docker Mount**: `./database/schemas:/docker-entrypoint-initdb.d`

## Benefits of Current Approach

- ✅ Simpler development workflow
- ✅ No schema drift between environments
- ✅ Consistent database state
- ✅ Faster iteration on schema changes
- ✅ No migration file maintenance overhead
