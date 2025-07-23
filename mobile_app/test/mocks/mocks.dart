import 'package:auto_light_pi/features/authentication/domain/use_cases/login_use_case.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:mocktail/mocktail.dart';

/* Domain Mocks */
class MockLoginUseCase extends Mock implements LoginUseCase {}

/* Data Mocks */

/* Presentation Mocks */
class MockAuthenticationBloc
    extends MockBloc<AuthenticationEvent, AuthenticationState>
    implements AuthenticationBloc {}

class MockLoginBloc extends MockBloc<LoginEvent, LoginState>
    implements LoginBloc {}
