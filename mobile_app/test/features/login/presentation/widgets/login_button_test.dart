import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:auto_light_pi/features/login/presentation/widgets/login_button.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:formz/formz.dart';
import 'package:mocktail/mocktail.dart';

import '../../../../mocks/mocks.dart';

void main() {
  late MockLoginBloc mockLoginBloc;

  setUp(() {
    mockLoginBloc = MockLoginBloc();
  });

  tearDown(() {
    mockLoginBloc.close();
  });

  group('LoginButton', () {
    testWidgets('renders correctly with default state', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(const LoginState());

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const LoginButton(),
            ),
          ),
        ),
      );

      expect(find.byType(ElevatedButton), findsOneWidget);
      expect(find.text('Submit'), findsOneWidget);
    });

    testWidgets('button is disabled when username/password are invalid', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(
        const LoginState(isUsernameValid: false, isPasswordValid: false),
      );

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const LoginButton(),
            ),
          ),
        ),
      );

      final ElevatedButton button = tester.widget<ElevatedButton>(
        find.byType(ElevatedButton),
      );

      expect(button.onPressed, isNull);
    });

    testWidgets('shows spinner when status is inProgress', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(
        const LoginState(
          isUsernameValid: true,
          isPasswordValid: true,
          status: FormzSubmissionStatus.inProgress,
        ),
      );

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const LoginButton(),
            ),
          ),
        ),
      );

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('dispatches LoginSubmitted event on button press', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(
        const LoginState(isUsernameValid: true, isPasswordValid: true),
      );

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const LoginButton(),
            ),
          ),
        ),
      );

      await tester.tap(find.byType(ElevatedButton));
      await tester.pumpAndSettle();

      verify(() => mockLoginBloc.add(const LoginSubmitted())).called(1);
    });
  });
}
