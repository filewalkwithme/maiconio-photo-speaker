# Photo-Speaker

Photo-Speaker is an application capable of identify labels in a photo and then speak these labels in english.

You can upload new photos and click in the button "Speak about!" to listen a description about each photo. (works only on Chrome and Firefox)

Photo-Speakear was built on top of IBM Bluemix Platform. It makes use of the Go runtime and other Bluemix services as Visual Recognition, Text to Speech and ElephantSQL(postgres)

## Try it alive!

http://maiconio-photo-speaker.mybluemix.net/

## Setting up Bluemix environment

1. Create a Bluemix account
2. On you dashboard, Click on 'Create APP' for Cloud Foundry.
3. Select 'WEB'
4. Choose the Go runtime
5. Choose a name for your APP
6. Click in 'Overview', then 'Add Service or API'
7. Add Visual Recognition, Text to Speech and ElephantSQL services

## Setting up the database

Photo-Speaker uses PostgreSQL to store photo information.
After add ElephantSQL, you will need to setup the following schema in your database:

```
CREATE TABLE photos
(
  id serial NOT NULL,
  file character varying NOT NULL,
  content text,
  CONSTRAINT photos_pkey PRIMARY KEY (id)
);

CREATE TABLE photo_labels
(
  id serial NOT NULL,
  photo_id integer NOT NULL,
  label character varying NOT NULL,
  score character varying NOT NULL,
  CONSTRAINT photo_labels_pkey PRIMARY KEY (id),
  CONSTRAINT photo_labels_photo_id_fkey FOREIGN KEY (photo_id)
      REFERENCES photos (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);
```

## Publishing Speaker-Photo

1. [Install Cloud Foundry command line tools][]
2. `cd` into in your local `speaker-photo` repository
3. then run:

```
cf api https://api.ng.bluemix.net
cf login -u <your-account-email> -o <your-account-email> -s dev
cf push <your-bluemix-app-name>
```

## Run the app locally

1. [Install Go][]
2. 'cd' into your speaker-photo repository
3. Run `go build maiconio-go`
4. Run `./maiconio-go`
5. Access the running app in a browser at http://localhost:8080

[Install Go]: https://golang.org/doc/install
[Install Cloud Foundry command line tools]: https://github.com/cloudfoundry/cli/releases
