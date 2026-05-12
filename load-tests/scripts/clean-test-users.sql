-- Удаляем всех пользователей, созданных нагрузочным тестом (username начинается с 'loadtest_')
DELETE FROM user_balances WHERE user_id IN (SELECT id FROM auth_credentials WHERE username LIKE 'loadtest_%');
DELETE FROM user_profiles WHERE user_id IN (SELECT id FROM auth_credentials WHERE username LIKE 'loadtest_%');
DELETE FROM audit_log WHERE user_id_hash IN (SELECT encode(sha256(cast(id as text)), 'hex') FROM auth_credentials WHERE username LIKE 'loadtest_%');
DELETE FROM documents WHERE user_id IN (SELECT id FROM auth_credentials WHERE username LIKE 'loadtest_%');
DELETE FROM analysis_results WHERE document_id IN (SELECT id FROM documents WHERE user_id IN (SELECT id FROM auth_credentials WHERE username LIKE 'loadtest_%'));
DELETE FROM auth_credentials WHERE username LIKE 'loadtest_%';
