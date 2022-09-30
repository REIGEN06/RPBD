## RPBD_1_Yakubov
# Выведите на экран любое сообщение
```sql
SELECT 'Hello I am Dan' as Любое_сообщение
```

# Выведите на экран текущую дату
```sql
SELECT CURRENT_DATE
```
# Создайте две числовые переменные и присвойте им значение. Выполните математические действия с этими числами и выведите результат на экран.
```sql
CREATE OR REPLACE FUNCTION dva_chisla() RETURNS int AS $$
DECLARE
	a int :=20;
	b int :=30;
BEGIN
	RETURN a+b;
END
$$ LANGUAGE plpgsql;
SELECT dva_chisla();
```

# Написать программу двумя способами 1 - использование IF, 2 - использование CASE. Объявите числовую переменную и присвоейте ей значение. Если число равно 5 - выведите на экран "Отлично". 4 - "Хорошо". 3 - Удовлетворительно". 2 - "Неуд". В остальных случаях выведите на экран сообщение, что введённая оценка не верна.
```sql
CREATE OR REPLACE FUNCTION sposobIF (a int) RETURNS char AS $$
BEGIN
	IF (a=5) THEN RETURN 'Отлично';
	ELSIF (a=4) THEN RETURN 'Хорошо';
	ELSIF (a=3) THEN RETURN 'Удовлетворительно';
	ELSIF (a=2) THEN RETURN 'Неуд';
 ELSE RETURN 'Введённая оценка не верна'
	END IF;
END
$$ LANGUAGE plpgsql;
```
```sql
CREATE OR REPLACE FUNCTION sposobCASE (a int) RETURNS char AS $$
BEGIN
	CASE WHEN (a=5) THEN RETURN 'Отлично';
	WHEN (a=4) THEN RETURN 'Хорошо';
	WHEN (a=3) THEN RETURN 'Удовлетворительно';
	WHEN (a=2) THEN RETURN 'Неуд';
 	WHEN (a>5 OR a<2) THEN RETURN 'Введённая оценка не верна';
	END CASE;
END
$$ LANGUAGE plpgsql;
SELECT sposobIF(2) AS ОценкаIF;
SELECT sposobCASE(1) AS ОценкаCASE;
```

# Выведите все квадраты чисел от 20 до 30 3-мя разными способами (LOOP, WHILE, FOR).
```sql
CREATE OR REPLACE PROCEDURE raiseLOOP() AS $$
DECLARE 
	x int := 20;
BEGIN
	LOOP
		RAISE NOTICE '%', x*x;
		x := x+1;
		EXIT WHEN x>30;
	END LOOP;
END
$$ LANGUAGE plpgsql;
CALL raiseLOOP();
```
```sql
CREATE OR REPLACE PROCEDURE raiseWHILE() AS $$
DECLARE 
	x int := 20;
BEGIN
	WHILE (x<31) LOOP
		RAISE NOTICE '%', x*x;
		x := x+1;
	END LOOP;
END
$$ LANGUAGE plpgsql;
CALL raiseWHILE();
```
```sql
CREATE OR REPLACE PROCEDURE raiseFOR() AS $$
DECLARE 
	x int := 20;
BEGIN
	FOR x in 20..30 LOOP
		RAISE NOTICE '%', x*x;
		x := x+1;
	END LOOP;
END
$$ LANGUAGE plpgsql;
CALL raiseFOR();
```

# Последовательность Коллатца. Берётся любое натуральное число. Если чётное - делим его на 2, если нечётное, то умножаем его на 3 и прибавляем 1. Такие действия выполняются до тех пор, пока не будет получена единица. Гипотеза заключается в том, что какое бы начальное число n не было выбрано, всегда получится 1 на каком-то шаге. Задания: написать функцию, входной параметр - начальное число, на выходе - количество чисел, пока не получим 1; написать процедуру, которая выводит все числа последовательности. Входной параметр - начальное число.
```sql
CREATE OR REPLACE FUNCTION COLLATSE (a int) RETURNS int AS $$
DECLARE
count int := 0;
BEGIN
	WHILE (a != 1)LOOP
		IF (a%2=0) THEN a := a/2;
		ELSE a := a*3+1;
		END IF;
		count:= count+1;
	END LOOP;
RETURN count;
END
$$ LANGUAGE plpgsql;
SELECT COLLATSE(10);
```
```sql
CREATE OR REPLACE PROCEDURE procedureCOLLATSE (a int) AS $$
BEGIN
	WHILE (a != 1)LOOP
		IF (a%2=0) THEN a := a/2;
		ELSE a := a*3+1;
		END IF;
		RAISE NOTICE '%', a;
	END LOOP;
END
$$ LANGUAGE plpgsql;
CALL procedureCOLLATSE(10);
```

# Числа Люка. Объявляем и присваиваем значение переменной - количество числе Люка. Вывести на экран последовательность чисел. Где L0 = 2, L1 = 1 ; Ln=Ln-1 + Ln-2 (сумма двух предыдущих чисел). Задания: написать фунцию, входной параметр - количество чисел, на выходе - последнее число (Например: входной 5, 2 1 3 4 7 - на выходе число 7); написать процедуру, которая выводит все числа последовательности. Входной параметр - количество чисел.

Напишите функцию, которая возвращает количество человек родившихся в заданном году.
Напишите функцию, которая возвращает количество человек с заданным цветом глаз.
Напишите функцию, которая возвращает ID самого молодого человека в таблице.
Напишите процедуру, которая возвращает людей с индексом массы тела больше заданного. ИМТ = масса в кг / (рост в м)^2.
Измените схему БД так, чтобы в БД можно было хранить родственные связи между людьми. Код должен быть представлен в виде транзакции (Например (добавление атрибута): BEGIN; ALTER TABLE people ADD COLUMN leg_size REAL; COMMIT;). Дополните БД данными.
Напишите процедуру, которая позволяет создать в БД нового человека с указанным родством.
Измените схему БД так, чтобы в БД можно было хранить время актуальности данных человека (выполнить также, как п.12).
Напишите процедуру, которая позволяет актуализировать рост и вес человека.
