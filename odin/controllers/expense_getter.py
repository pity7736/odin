from odin.repositories import ExpenseRepository


class ExpenseGetter:

    def get_by_uuid(self, uuid):
        repository = ExpenseRepository()
        try:
            return repository.get_by(uuid=uuid)
        except KeyError:
            return None
