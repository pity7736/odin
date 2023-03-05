CREATE MIGRATION m1u64rhbzpr37os7dqk3avbbwsyr5bqnpdbntappqpeziusnfsfsrq
    ONTO initial
{
  CREATE TYPE default::Category {
      CREATE REQUIRED PROPERTY name -> std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE default::Wallet {
      CREATE REQUIRED PROPERTY balance -> std::decimal;
      CREATE REQUIRED PROPERTY name -> std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE default::Movement {
      CREATE REQUIRED LINK category -> default::Category;
      CREATE REQUIRED LINK wallet -> default::Wallet;
      CREATE REQUIRED PROPERTY amount -> std::decimal;
      CREATE REQUIRED PROPERTY date -> cal::local_date;
      CREATE REQUIRED PROPERTY type -> std::str {
          CREATE CONSTRAINT std::one_of('expense', 'income');
      };
  };
  CREATE TYPE default::Transfer {
      CREATE REQUIRED LINK expense -> default::Movement;
      CREATE REQUIRED LINK income -> default::Movement;
      CREATE REQUIRED LINK source -> default::Wallet;
      CREATE REQUIRED LINK target -> default::Wallet;
      CREATE REQUIRED PROPERTY amount -> std::decimal;
      CREATE REQUIRED PROPERTY date -> cal::local_date;
  };
  CREATE TYPE default::User {
      CREATE REQUIRED PROPERTY email -> std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY first_name -> std::str;
      CREATE REQUIRED PROPERTY last_name -> std::str;
      CREATE REQUIRED PROPERTY password -> std::str;
  };
  CREATE TYPE default::Token {
      CREATE REQUIRED LINK user -> default::User;
      CREATE REQUIRED PROPERTY value -> std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
};
