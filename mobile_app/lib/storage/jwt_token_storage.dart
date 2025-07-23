import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class JwtTokenStorage {
  static const String _jwtKey = 'jwt';
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  Future<void> writeToken(String token) async {
    await _storage.write(key: _jwtKey, value: token);
  }

  Future<String?> readToken() async {
    return _storage.read(key: _jwtKey);
  }

  Future<void> deleteToken() async {
    await _storage.delete(key: _jwtKey);
  }
}
