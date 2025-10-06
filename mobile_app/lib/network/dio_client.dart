import 'package:auto_light_pi/interceptors/jwt_token_interceptor.dart';
import 'package:dio/dio.dart';

class DioClient {
  final Dio dio;

  DioClient._(this.dio);

  static Future<DioClient> create({
    required String baseUrl,
    required JwtTokenInterceptor jwtTokenInterceptor,
  }) async {
    final Dio dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: const Duration(milliseconds: 5000),
        receiveTimeout: const Duration(milliseconds: 5000),
      ),
    );

    dio.interceptors.add(jwtTokenInterceptor);

    return DioClient._(dio);
  }
}
