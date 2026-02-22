import 'package:auto_light_pi/features/authentication/data/data_sources/auth_local_data_source.dart';
import 'package:auto_light_pi/features/authentication/data/data_sources/auth_remote_data_source.dart';
import 'package:auto_light_pi/features/authentication/data/repositories/auth_repository_impl.dart';
import 'package:auto_light_pi/features/authentication/domain/repositories/auth_repository.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/check_authentication.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/login.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/logout.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/network/dio_client.dart';
import 'package:auto_light_pi/interceptors/jwt_token_interceptor.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:get_it/get_it.dart';

final GetIt sl = GetIt.instance;

void init() async {
  /** Dio Client **/
  sl.registerSingletonAsync<DioClient>(
    () async => await DioClient.create(
      baseUrl: dotenv.get('BACKEND_URL', fallback: 'http://localhost:8080/api'),
      jwtTokenInterceptor: sl<JwtTokenInterceptor>(),
    ),
  );

  /** Interceptors **/
  sl.registerLazySingleton(
    () => JwtTokenInterceptor(sl<AuthLocalDataSource>()),
  );

  /** Domain **/
  sl.registerLazySingleton(() => LoginUseCase(sl()));
  sl.registerLazySingleton(() => CheckAuthenticationUseCase(sl()));
  sl.registerLazySingleton(() => LogoutUseCase(sl()));
  sl.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(
      remote: sl<AuthRemoteDataSource>(),
      local: sl<AuthLocalDataSource>(),
    ),
  );

  /** Data **/
  sl.registerLazySingleton(() => AuthRemoteDataSource(sl<DioClient>()));
  sl.registerLazySingleton(() => AuthLocalDataSource());

  /** Presentation **/
  sl.registerFactory(() => LoginBloc(authenticationBloc: sl()));
  sl.registerLazySingleton(
    () => AuthenticationBloc(
      loginUseCase: sl(),
      checkAuthenticationUseCase: sl(),
      logoutUseCase: sl(),
    ),
  );

  await sl.allReady();
}
