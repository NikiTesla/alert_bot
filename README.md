# Alert bot
Inititally made for receiving results of computer modelling or other long-lasting calculations and programs

### Well, it is not complicated for now and it has:
- Service that receives both text results on `/result` handler and images on `/image` handler
- Telegram package that is responsible for interaction with telegram bot (you should create it with Botfather and specify your own `TELEBOT_API` env) 
- Storage with for now only subscriptions info (which users should get messages with updates: either text or images)
