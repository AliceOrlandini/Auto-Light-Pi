import 'package:auto_light_pi/features/authentication/data/models/user.dart';
import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'dart:convert';

class AuthLocalDataSource {
  final FlutterSecureStorage _secureStorage = const FlutterSecureStorage();
  final String _userKey = 'user';
  final String _jwtKey = 'jwt';

  Future<void> cacheUser(User user) async {
    final String userJson = jsonEncode(user.toJson());
    await _secureStorage.write(key: _userKey, value: userJson);
  }

  Future<void> cacheToken(String token) async {
    await _secureStorage.write(key: _jwtKey, value: token);
  }

  Future<UserEntity?> getUser() async {
    final String? userJson = await _secureStorage.read(key: _userKey);
    if (userJson == null) return null;

    return User.fromJson(jsonDecode(userJson)).toEntity();
  }

  Future<String?> getToken() async {
    return await _secureStorage.read(key: _jwtKey);
  }

  Future<void> clearUser() async {
    await _secureStorage.delete(key: _userKey);
  }

  Future<void> clearToken() async {
    await _secureStorage.delete(key: _jwtKey);
  }
}
