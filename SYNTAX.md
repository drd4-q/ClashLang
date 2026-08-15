# ClashLang — Справочник синтаксиса

## Запуск

```bash
./clashlang file.clash            # выполнение
./clashlang file.clash --debug    # выполнение + отладка после
```

## Основы

Команды **без регистра**. Имена переменных **с регистром**.

```clash
function(MyFunc){   // OK — функция MyFunc
    c = 12          // OK — переменная c
    C = 12          // это ДРУГАЯ переменная C
}
```

Комментарии: `// текст до конца строки`

---

## Переменные

```clash
c = 12             // int
pi = 3.14          // float
name = 'Alex'      // string
flag = true        // bool
```

---

## Математика

```clash
solve(a + b)       // результат в lastResult
solve.out = var    // сохранить результат
```

Операторы: `+`, `-`, `*`, `/`, `^`

---

## Текст

```clash
text(a + b)        // конкатенация
text.out = var     // сохранить результат
```

---

## Ввод

```clash
solve.input() = var    // ввод числа
text.input() = var     // ввод текста
input()                // общий ввод (авто-тип)
solve.out = var        // сохранить результат ввода
```

---

## Печать

```clash
print(var)         // печать переменной или литерала
memory.out         // печать всех переменных (зашифровано)
```

---

## Функции

```clash
function(Name){ тело }       // определение с телом
function(Name)               // определение имени
memory.start{ тело }         // тело отдельно

jump(Name)                   // вызов
memory.start(Name)           // альтернативный вызов
memory.load(A, B)            // вызов нескольких
pass()                       // пропуск
```

---

## Условия

```clash
if c = 12 then jump(Print)       // однострочный
if c = 12                        // двустрочный
then jump(Print)
if c = 12 then jump(Yes)         // с else
else: jump(No)
if file.none = jump(Missing)     // файл не найден
```

---

## Файлы

```clash
file.open(name)
file.write(text)
file.close
file.text.read
file.text.copy = var
file.create(name)
```

---

## Циклы

```clash
repeat(3){               // повторить 3 раза
    print(i)
}

i = 0
repeat(i){               // повторить i раз
    solve(i + 1)
    solve.out = i
}
```

---

## Строковые операции

```clash
len(var)                    // длина строки → lastResult
len.out = var               // сохранить

upper(var)                  // в верхний регистр
upper.out = var

lower(var)                  // в нижний регистр
lower.out = var

reverse(var)                // реверс строки
reverse.out = var

replace(var, old, new)      // замена подстроки
replace.out = var

substr(var, start, len)     // подстрока
substr.out = var

find(var, sub)              // позиция подстроки (-1 = нет)
find.out = var
```

---

## Математика (доп.)

```clash
abs(var)                    // абсолютное значение
abs.out = var

random()                    // случайное число 0-99
random.out = var
```

---

## Типы

```clash
type(var)                   // тип переменной: int, float, string, bool
type.out = var
```

---

## Отладка

```clash
debug()                     // все переменные и функции (открытый вид)
decode()                    // расшифровка memory.out
assert(var, value)          // проверка значения
clear()                     // очистить все переменные
```

---

## Защита памяти

```clash
memory.out      // зашифрованный вывод (XOR + hex)
debug()         // открытый вывод
decode()        // расшифрованный вывод
```

---

## Пример

```clash
function(Hello){
    name = 'John'
    age = 25

    upper(name)
    upper.out = nameUpper
    print(nameUpper)

    len(name)
    len.out = nameLen
    print(nameLen)

    repeat(2){
        print(age)
    }

    assert(age, 25)
    debug()
}

memory.start(Hello)
memory.out
decode()
```
