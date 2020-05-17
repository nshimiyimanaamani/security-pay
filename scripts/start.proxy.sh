export TOKEN=$(head -c 16 /dev/urandom | shasum | cut -f1 -d" ")
echo "Token is: ${TOKEN}"

export HEROKU_APP=codechef-inlets

heroku config:set TOKEN=${TOKEN} --app $HEROKU_APP
echo $TOKEN > token


inlets client \
  --remote wss://${HEROKU_APP}.herokuapp.com \
  --token $(cat token) \
  --upstream http://127.0.0.1:8000