@startuml
skinparam backgroundColor #FEFECE
skinparam handwritten true

title Диаграмма последовательностей для HTTP API системы расчёта баллов лояльности

participant "Пользователь" as User
participant "Система лояльности" as LoyaltySystem
participant "Система расчёта баллов" as PointsCalcSystem

note right of User:Регистрация и аутентификация пользователя
User -> LoyaltySystem:POST /api/user/register — регистрация пользователя
User -> LoyaltySystem:POST /api/user/login — аутентификация пользователя

note right of User:Работа с системой лояльности
User -> LoyaltySystem:POST /api/user/orders — загрузка номера заказа для расчёта

note right of LoyaltySystem:Работа с системой расчета балов
LoyaltySystem --> PointsCalcSystem:POST /api/orders — регистрация нового заказа
LoyaltySystem --> PointsCalcSystem:GET /api/orders/{number} — информации о расчёте 

User -> LoyaltySystem:GET /api/user/orders — список загруженных номеров заказов
User -> LoyaltySystem:GET /api/user/balance — текущий баланс счёта баллов
User -> LoyaltySystem:POST /api/user/balance/withdraw — списание баллов
User -> LoyaltySystem:GET /api/user/withdrawals — информация о списании средств

note right of User:Регистрация логики начилсения балов
User --> PointsCalcSystem:POST /api/goods — регистрация информации о новой механике вознаграждения за товар.
@enduml


accrual

gophermart
    //регистрация
        /login
        /registration
    //заказы
        /загрузка пользователем номера заказа для расчета балов - LoadOrder
        /получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях ListOrders
    //работа с балансом балов
        /получение текущего баланса счета балов лояльности пользователя
        /запрос на списание балов с накопительного счета в счет оплаты нового заказа
        /получение информации о выводе средств с накопительного счета пользователем


mockgen -source=internal/gophermart/services/user/user_service.go -destination=internal/gophermart/services/user/mock_user_service.go -package=user
mockgen -destination internal/gophermart/services/user/mock_pgx_tx.go -package mocks github.com/jackc/pgx/v5 Tx



