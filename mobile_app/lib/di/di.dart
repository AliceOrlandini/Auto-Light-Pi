import 'package:auto_light_pi/features/authentication/data/data_sources/auth_remote_data_source.dart';
import 'package:auto_light_pi/features/authentication/data/repositories/auth_repository_impl.dart';
import 'package:auto_light_pi/features/authentication/domain/repositories/auth_repository.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/login_use_case.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/network/dio_client.dart';
import 'package:auto_light_pi/storage/jwt_token_storage.dart';
import 'package:get_it/get_it.dart';

final GetIt sl = GetIt.instance;

void init() async {
  /** Local Storage **/
  sl.registerLazySingleton<JwtTokenStorage>(() => JwtTokenStorage());

  /** Dio Client **/
  sl.registerSingletonAsync<DioClient>(
    () async => await DioClient.create(
      baseUrl: 'http://192.168.1.5:8080/api',
      jwtTokenStorage: sl<JwtTokenStorage>(),
    ),
  );

  /** Domain **/
  sl.registerLazySingleton(() => LoginUseCase(sl()));
  sl.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(remote: sl<AuthRemoteDataSource>()),
  );

  /** Data **/
  sl.registerLazySingleton(() => AuthRemoteDataSource(sl<DioClient>()));

  /** Presentation **/
  sl.registerFactory(() => LoginBloc(authenticationBloc: sl()));
  sl.registerFactory(() => AuthenticationBloc(loginUseCase: sl()));

  await sl.allReady();
}
