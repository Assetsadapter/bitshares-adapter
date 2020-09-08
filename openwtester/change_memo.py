from bitshares.transactionbuilder import TransactionBuilder
from bitsharesbase.operations import Account_update
from bitshares import BitShares
from bitshares import storage
from bitshares.wallet import Wallet
from bitsharesbase.account import PrivateKey
# wifs = [
#     "5KDDijbYFxmWNTdgqvKhQPsnDrdszrqt5kWYk5p8JijqdPbpRJE"
# ]
# wif1=PrivateKey("5KDDijbYFxmWNTdgqvKhQPsnDrdszrqt5kWYk5p8JijqdPbpRJE")
#账户对应wif私钥
wif2=PrivateKey("5KDDijbYFxmWNTdgqvKhQPsnDrdszrqt5kWYk5p8JijqdPbpRJE")
config = storage.InRamConfigurationStore()

key_store = storage.InRamPlainKeyStore(config=config)
wallet = Wallet(key_store=key_store)

b = BitShares(node="http://3.wallet.info/bts",wallet=wallet)

# b.wallet = wallet

# b.unlock("iexbts666")
# wallet.addPrivateKey(wif1)
wallet.addPrivateKey(wif2)
tx=b.update_memo_key("BTS6GfoTwr3VY3scZbKK7XmsEtYyvso6M3buPGbbpJYhVTb6uaLn8","1.2.1791970")

# tx = TransactionBuilder(bitshares_instance=b)
#
# op = Account_update(
#         **{
#             "fee": {"amount": 1467634, "asset_id": "1.3.0"},
#             "account": "1.2.18",
#             "new_options": {
#                 "memo_key": "BTS6uCBue3rwdPHzULtpSSHRExwfKtA5xUMSm27s2AipsvSLf8V1B",
#                 "voting_account": "1.2.5",
#                 "num_witness": 0,
#                 "num_committee": 0,
#                 "votes": [],
#                 "extensions": [],
#
#             }
#         }
#     )
# tx.appendOps(op)
# tx.appendWif("5JvHXy6F2gr5WAAZjmmxVecrtaMYYKCwdJu7LV3mkQKrJZwWpFW")
# tx.sign()
# print(tx.json())
# tx.broadcast()

