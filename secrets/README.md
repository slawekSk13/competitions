# Secrets Directory

This directory contains sensitive configuration files for the application. These files should never be committed to version control.

## Files

- `db_name.txt`: Database name
- `db_user.txt`: Database username
- `db_password.txt`: Database password
- `pgadmin_email.txt`: pgAdmin email address
- `pgadmin_password.txt`: pgAdmin password

## Security Notice

These files contain sensitive information and should be:

1. Added to .gitignore
2. Kept secure and not shared
3. Backed up securely
4. Used only in development and production environments

## Security Best Practices

- Use strong, unique passwords
- Set proper file permissions: `chmod 600 *.txt`
