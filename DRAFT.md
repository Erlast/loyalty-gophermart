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


