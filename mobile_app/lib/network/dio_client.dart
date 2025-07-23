import 'package:auto_light_pi/interceptors/jwt_token_interceptor.dart';
import 'package:auto_light_pi/storage/jwt_token_storage.dart';
import 'package:dio/dio.dart';

class DioClient {
  final Dio dio;

  DioClient._(this.dio);

  static Future<DioClient> create({
    required String baseUrl,
    required JwtTokenStorage jwtTokenStorage,
  }) async {
    final Dio dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: const Duration(milliseconds: 5000),
        receiveTimeout: const Duration(milliseconds: 5000),
      ),
    );
    dio.interceptors.add(JwtTokenInterceptor(jwtTokenStorage));
    return DioClient._(dio);
  }
}
