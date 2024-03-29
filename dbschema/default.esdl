module default {
    type Wallet {
        required property balance -> decimal;
        required property name -> str {
            constraint exclusive;
        }
    }

    type Category {
        required property name -> str {
            constraint exclusive;
        }
    }

    type Movement {
        required property date -> cal::local_date;
        required property amount -> decimal;
        required link category -> Category;
        required link wallet -> Wallet;
        required property type -> str {
            constraint one_of ('expense', 'income');
        }
    }

    type Transfer {
        required link source -> Wallet;
        required link target -> Wallet;
        required link expense -> Movement;
        required link income -> Movement;
        required property amount -> decimal;
        required property date -> cal::local_date;
    }

    type User {
        required property email -> str {
            constraint exclusive;
        }
        required property password -> str;
        required property first_name -> str;
        required property last_name -> str;
    }

    type Token {
        required link user -> User;
        required property value -> str {
            constraint exclusive;
        }
    }
};
