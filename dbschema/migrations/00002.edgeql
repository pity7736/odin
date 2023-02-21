CREATE MIGRATION m1vulsj76qvzv6wsog34aeakdhex3t42r4qodcerkj3nt7bve3mdoa
    ONTO m1yaikwyzm37ooaeudh6topuybl2pj6bko5anwu2m64cdypq2totdq
{
  ALTER TYPE default::Token {
      ALTER PROPERTY value {
          CREATE CONSTRAINT std::exclusive;
      };
  };
};
