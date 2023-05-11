# TASK #4

---

# **Tables**

## PLAYERS
| Атрибут       | Тип         | Описание                         | Значение |
|---------------|-------------|----------------------------------|----------|
| id     *      | int(serial) | *primary key*                    | 1...n    |
| progress_id * | int         | *foreign key(DAILYBOX_PROGRESS)* |          |
| email         | string      |                                  |          |
| pwd_hash      | string      | *encrypted data*                 |          |
| money         | int         |                                  |          | 
| crystals      | int         |                                  |          |

## DAILYBOX_REWARDS
| Атрибут      | Тип         | Описание                | Значение    |
|--------------|-------------|-------------------------|-------------|
| id     *     | int(serial) | *id*                    | 1...lastday |
| progres_id * | int         | *foreign key(DAILYBOX)* |             |
| money        | int         |                         |             |
| crystals     | int         |                         |             |

## DAILYBOX_PROGRESS
| Атрибут     | Тип         | Описание                        | Значение |
|-------------|-------------|---------------------------------|----------|
| id     *    | int(serial) | *id*                            | 1...n    |
| user_id *   | int         | *foreign key(PLAYERS)*          |          |
| reward_id * | int         | *foreign key(DAILYBOX_REWARDS)* |          |
| rewarded *  | bool        |                                 |          |
| last_update | timstamp    |                                 |          |

---

## **Шаблон метода на выдачу награды**

### POST v1/currGameId:int/reward

**request**

```http request
TYPE: POST
HEADER: autorization=accessToken
URL: empty

Payload: {
"id":1
}
```

**response**
```http request
{
   "id":1,
   "rewarded":true,
}
```

---
## **Шаблон структуры логики**
 1. Валидируем токен игрока (токен не валиден => errInvalidToken)
 2. Начинаем транзакцию (уровень изоляции - for update)
 3. Проверям получал ли пользователь в этот день награду или является ли текущий день награды последним (получал или является => закрываем транзакцию и отдаем errAlreadyRewarded)
 4. Иначе, обновляем баланс игрока (user.money = user.money + daylibox_reward.money && user.crystals = user.crystals + daylibox_reward.crystals). Ставим rewarded true и перемещаемся на инкрементируемся на следующий день награды (при условии, что curr.dailybox_reward_id < lastday) 
 5. Закрываем транзакцию и отдаем response