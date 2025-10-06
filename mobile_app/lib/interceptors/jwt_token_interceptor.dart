import 'package:auto_light_pi/features/authentication/data/data_sources/auth_local_data_source.dart';
import 'package:dio/dio.dart';

class JwtTokenInterceptor extends Interceptor {
  final AuthLocalDataSource _local;

  JwtTokenInterceptor(this._local);

  @override
  void onResponse(
    Response<dynamic> response,
    ResponseInterceptorHandler handler,
  ) async {
    final List<String>? setCookie = response.headers.map['set-cookie'];
    if (setCookie != null) {
      final RegExpMatch? match = RegExp(
        r'jwt=([^;]+)',
      ).firstMatch(setCookie.join(';'));
      if (match != null) {
        await _local.cacheToken(match.group(1)!);
      }
    }
    handler.next(response);
  }

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final String? token = await _local.getToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }
}
