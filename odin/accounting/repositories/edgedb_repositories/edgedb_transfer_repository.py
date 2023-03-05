from odin.accounting.models import Transfer, Wallet
from .db_client import DBClient
from ..repositories import TransferRepository


class EdgeDBTransferRepository(TransferRepository):

    def __init__(self):
        self._client = DBClient()

    def add(self, transfer):
        source = 'select Wallet filter .name = <str>$source_name'
        target = 'select Wallet filter .name = <str>$target_name'
        expense = 'select Movement filter .id = <uuid>$expense_uuid'
        income = 'select Movement filter .id = <uuid>$income_uuid'
        query = (
            'insert Transfer {'
            'amount := <decimal>$amount, date := <cal::local_date>$date, '
            f'source := ({source}), target := ({target}), expense := ({expense}), income := ({income})}}'
        )
        result = self._client.query_single(
            query,
            source_name=transfer.source.name,
            target_name=transfer.target.name,
            expense_uuid=transfer.expense.uuid,
            income_uuid=transfer.income.uuid,
            amount=transfer.amount,
            date=transfer.date
        )
        transfer.uuid = result.id

    def get_by_uuid(self, uuid: str):
        record = self._client.query_single(
            'select Transfer {amount, date, source: {name, balance}, target: {name, balance}} filter .id = <uuid>$uuid',
            uuid=uuid
        )
        if record:
            return Transfer(
                amount=record.amount,
                date=record.date,
                source=Wallet(name=record.source.name, balance=record.source.balance),
                target=Wallet(name=record.target.name, balance=record.target.balance)
            )
