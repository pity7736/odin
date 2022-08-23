from odin.accounting.models import Income


class IncomeRepository:

    _incomes = {}

    def add(self, income: Income):
        self.__class__._incomes[income.uuid] = income

    def get_by_uuid(self, uuid: str) -> Income:
        return self._incomes.get(uuid)
