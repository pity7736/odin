CREATE MIGRATION m1yaikwyzm37ooaeudh6topuybl2pj6bko5anwu2m64cdypq2totdq
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
  CREATE ABSTRACT TYPE default::Movement {
      CREATE REQUIRED LINK category -> default::Category;
      CREATE REQUIRED LINK wallet -> default::Wallet;
      CREATE REQUIRED PROPERTY amount -> std::decimal;
      CREATE REQUIRED PROPERTY date -> cal::local_date;
  };
  CREATE TYPE default::Expense EXTENDING default::Movement;
  CREATE TYPE default::Income EXTENDING default::Movement;
  CREATE TYPE default::Transference {
      CREATE REQUIRED LINK expense -> default::Expense;
      CREATE REQUIRED LINK income -> default::Income;
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
      CREATE REQUIRED PROPERTY value -> std::str;
  };
};
