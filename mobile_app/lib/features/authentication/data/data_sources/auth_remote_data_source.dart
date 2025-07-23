import 'package:auto_light_pi/features/authentication/data/models/login_response.dart';
import 'package:auto_light_pi/network/dio_client.dart';
import 'package:auto_light_pi/network/network_exception.dart';
import 'package:dio/dio.dart';

class AuthRemoteDataSource {
  final DioClient _dioClient;

  AuthRemoteDataSource(this._dioClient);

  Future<LoginResponse> loginByUsername(
    String username,
    String password,
  ) async {
    try {
      final Response<dynamic> response = await _dioClient.dio.post(
        '/login/username',
        data: <String, String>{'username': username, 'password': password},
      );
      return LoginResponse.fromJson(response.data as Map<String, dynamic>);
    } on DioException catch (e) {
      if (e.response != null && e.response!.data is Map<String, dynamic>) {
        final Map<String, dynamic> data =
            e.response!.data as Map<String, dynamic>;
        final String errorMsg = data['error'] as String? ?? 'unknown error';

        switch (e.response!.statusCode) {
          case 400:
            throw BadRequestException(errorMsg);
          case 401:
            throw UnauthorizedException(errorMsg);
          case 408:
            throw TimeoutException(errorMsg);
          case 500:
            throw ServerException(errorMsg);
          default:
            throw NetworkException(
              errorMsg,
              statusCode: e.response!.statusCode,
            );
        }
      }
      if (e.message != null) {
        throw NetworkException(e.message!);
      }
      // if no response or exception was thrown, throw a generic exception
      throw Exception('unexpected error occurred during login');
    }
  }
}
