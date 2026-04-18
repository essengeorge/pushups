while true; do
    env $(cat .env | xargs) ./punkpushups
    ./punkpushups
    sleep 5
done
