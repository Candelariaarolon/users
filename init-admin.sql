-- Script para crear el primer usuario administrador
-- Ejecutar después de que la base de datos esté inicializada

-- CREDENCIALES DEL ADMIN:
-- Email:    admin@unichat.com
-- Password: admin123
-- Nota: Cambiar la contraseña después del primer login por seguridad

INSERT INTO user_models (email, password_hash, first_name, last_name, is_admin, is_verified, created_at)
SELECT 'admin@unichat.com',
       '240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9', -- Hash SHA-256 de "admin123"
       'Admin',
       'UniChat',
       true,
       true,
       NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM user_models WHERE email = 'admin@unichat.com'
);
