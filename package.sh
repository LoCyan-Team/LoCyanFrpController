#!/bin/sh
set -e

# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

controller_version="0.0.4"
echo "build version: $controller_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64 loong64'
extra_all='_ hf'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        for extra in $extra_all; do
            suffix="${os}_${arch}"
            if [ "x${extra}" != x"_" ]; then
                suffix="${os}_${arch}_${extra}"
            fi
            controller_dir_name="controller_${controller_version}_${suffix}"
            controller_path="./packages/controller_${controller_version}_${suffix}"

            if [ "x${os}" = x"windows" ]; then
                if [ ! -f "./controller_${os}_${arch}.exe" ]; then
                    continue
                fi
                mkdir ${controller_path}
                mv ./controller_${os}_${arch}.exe ${controller_path}/controller.exe
            else
                if [ ! -f "./controller_${suffix}" ]; then
                    continue
                fi
                mkdir ${controller_path}
                mv ./controller_${suffix} ${controller_path}/controller
            fi  
            cp ../LICENSE ${controller_path}
            cp -f ../config.ini.example ${controller_path}

            # packages
            cd ./packages
            if [ "x${os}" = x"windows" ]; then
                zip -rq ${controller_dir_name}.zip ${controller_dir_name}
            else
                tar -zcf ${controller_dir_name}.tar.gz ${controller_dir_name}
            fi  
            cd ..
            rm -rf ${controller_path}
        done
    done
done

cd -
