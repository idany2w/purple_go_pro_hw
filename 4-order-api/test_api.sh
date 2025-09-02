#!/bin/bash

# Скрипт для тестирования API заказов
# Убедитесь, что сервер запущен на localhost:8080

BASE_URL="http://localhost:8080"
JWT_TOKEN=""

echo "🧪 Тестирование API заказов"
echo "================================"

# Функция для вывода результата
print_result() {
    if [ $1 -eq 0 ]; then
        echo "✅ $2"
    else
        echo "❌ $2"
    fi
}

# 1. Регистрация пользователя
echo "1. Регистрация пользователя..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"phone": "+79001234567"}')

if echo "$REGISTER_RESPONSE" | grep -q "SMS sent"; then
    print_result 0 "Регистрация успешна"
else
    print_result 1 "Ошибка регистрации: $REGISTER_RESPONSE"
fi

echo ""

# 2. Получение списка продуктов
echo "2. Получение списка продуктов..."
PRODUCTS_RESPONSE=$(curl -s "$BASE_URL/products")

if echo "$PRODUCTS_RESPONSE" | grep -q "products"; then
    print_result 0 "Продукты получены успешно"
    echo "   Продукты: $PRODUCTS_RESPONSE"
else
    print_result 1 "Ошибка получения продуктов: $PRODUCTS_RESPONSE"
fi

echo ""

# 3. Получение заказа по ID (должен вернуть 404, так как заказов еще нет)
echo "3. Попытка получить несуществующий заказ..."
ORDER_RESPONSE=$(curl -s -w "%{http_code}" "$BASE_URL/order/999")

HTTP_CODE="${ORDER_RESPONSE: -3}"
RESPONSE_BODY="${ORDER_RESPONSE%???}"

if [ "$HTTP_CODE" = "404" ]; then
    print_result 0 "Правильно возвращается 404 для несуществующего заказа"
else
    print_result 1 "Неожиданный HTTP код: $HTTP_CODE, ответ: $RESPONSE_BODY"
fi

echo ""

# 4. Попытка создать заказ без токена (должна вернуть 401)
echo "4. Попытка создать заказ без токена..."
CREATE_ORDER_RESPONSE=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/order" \
    -H "Content-Type: application/json" \
    -d '{"product_ids": [1], "quantities": [1]}')

HTTP_CODE="${CREATE_ORDER_RESPONSE: -3}"
RESPONSE_BODY="${CREATE_ORDER_RESPONSE%???}"

if [ "$HTTP_CODE" = "401" ]; then
    print_result 0 "Правильно возвращается 401 для неавторизованного запроса"
else
    print_result 1 "Неожиданный HTTP код: $HTTP_CODE, ответ: $RESPONSE_BODY"
fi

echo ""

echo "🎯 Тестирование завершено!"
echo ""
echo "📝 Примечания:"
echo "- Для полного тестирования с JWT токеном нужно сначала пройти аутентификацию"
echo "- Сервер должен быть запущен на $BASE_URL"
echo "- База данных должна быть доступна"
echo ""
echo "🚀 Для запуска сервера используйте: go run cmd/main.go"
