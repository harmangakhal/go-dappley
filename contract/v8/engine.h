#include <stdint.h>
#include <stdbool.h>
#include "lib/transaction_struct.h"
#include "lib/utxo_struct.h"

#ifdef WIN32 
#ifdef V8DLL
#define EXPORT __declspec(dllexport)
#else 
#define EXPORT __declspec(dllimport)
#endif
#else
#define EXPORT __attribute__((__visibility__("default")))
#endif 

#ifdef __cplusplus
extern "C" {
#endif
    typedef bool (*FuncVerifyAddress)(const char *address);
    typedef int (*FuncTransfer)(void *handler, const char *to, const char *amount, const char *tip);
    typedef char* (*FuncStorageGet)(void *address, const char *key);
    typedef int (*FuncStorageSet)(void *address, const char *key, const char *value);
    typedef int (*FuncStorageDel)(void *address, const char *key);
    typedef int (*FuncTriggerEvent)(void *address, const char *topic, const char *data);
    typedef void (*FuncTransactionGet)(void* address, void* context);
    typedef void (*FuncPrevUtxoGet)(void* address, void* context);
    typedef void (*FuncLogger)(unsigned int level, char** args, int length);
    typedef int (*FuncRecordReward)(void *handler, const char *address, const char *amount);
    typedef bool (*FuncVerifySignature)(const char *msg, const char *pubKey, const char *sig);
    typedef bool (*FuncVerifyPublicKey)(const char *addr, const char *pubKey);
    typedef int (*FuncRandom)(void *handler, int max);
    typedef int (*FuncGetCurrBlockHeight)(void *handler);
    typedef char* (*FuncGetNodeAddress)(void *handler);
	typedef void* (*FuncMalloc)(size_t size);
	typedef void  (*FuncFree)(void* data);

    EXPORT void Initialize();
    EXPORT int executeV8Script(const char *sourceCode, uintptr_t handler, char **result);
    EXPORT void InitializeBlockchain(FuncVerifyAddress verifyAddress, FuncTransfer transfer, FuncGetCurrBlockHeight getCurrBlockHeight, FuncGetNodeAddress getNodeAddress);
    EXPORT void InitializeRewardDistributor(FuncRecordReward recordReward);
    EXPORT void InitializeStorage(FuncStorageGet get, FuncStorageSet set, FuncStorageDel del);
    EXPORT void InitializeEvent(FuncTriggerEvent triggerEvent);
    EXPORT void InitializeTransaction(FuncTransactionGet get);
    EXPORT void InitializeCrypto(FuncVerifySignature verifySignature, FuncVerifyPublicKey verifyPublicKey);
    EXPORT void InitializeMath(FuncRandom random);
    EXPORT void SetTransactionData(struct transaction_t* tx, void* context);
    EXPORT void InitializePrevUtxo(FuncPrevUtxoGet get);
    EXPORT void SetPrevUtxoData(struct utxo_t* utxos, int length, void* context);
    EXPORT void InitializeLogger(FuncLogger logger);
    EXPORT void InitializeSmartContract(char* source);
    EXPORT void DisposeV8();
	EXPORT void InitializeMemoryFunc(FuncMalloc mallocFunc, FuncFree freeFunc);
#ifdef __cplusplus
}
#endif
