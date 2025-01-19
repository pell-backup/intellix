#! /bin/bash

buf generate

# avsi proto
mv generate/cosmos/intellix/dvs/price/types/* ../dvs/price/types/
mv generate/cosmos/intellix/dvs/processor/types/* ../dvs/processor/types/

# dvs proto
mv generate/cosmos/intellix/pkg/pelldvs/types/* ../pkg/pelldvs/types/

# cosmos proto
mv generate/cosmos/intellix/x/price/types/* ../x/price/types/
mv generate/cosmos/intellix/x/processor/types/* ../x/processor/types/

# pulsar proto
# mv generate/go_grpc/intellix/pelldvs/* ../api/intellix/pelldvs/
# mv generate/go_grpc/intellix/price/* ../api/intellix/price/
# mv generate/go_grpc/intellix/processor/* ../api/intellix/processor/

rm -rf generate