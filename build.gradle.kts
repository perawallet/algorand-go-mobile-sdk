plugins {
    alias(libs.plugins.android.library) apply false
    alias(libs.plugins.kotlin.android) apply false
    alias(libs.plugins.jetbrains.kotlin.jvm)
    id("maven-publish")
    `java-library`
}

allprojects {
    tasks.withType<JavaCompile> {
        options.encoding = "UTF-8"
    }
}

group = "com.algorand.algosdk"
version = "1.0.0"

publishing {
    publications {
        create<MavenPublication>("release") {
            groupId = group.toString()
            artifactId = "algosdk-android"
            version = version.toString()

            artifact("$projectDir/output/algosdk.aar")

            pom {
                name.set("Algorand Go Mobile SDK (Android)")
                description.set("Android AAR generated from the Go Mobile Algorand SDK.")
                url.set("https://github.com/algorand/go-mobile-algorand-sdk")

                licenses {
                    license {
                        name.set("Apache License 2.0")
                        url.set("https://www.apache.org/licenses/LICENSE-2.0")
                    }
                }

                developers {
                    developer {
                        id.set("AlgorandFoundation")
                        name.set("Algorand Foundation")
                        email.set("press@algorand.foundation")
                    }
                }

                scm {
                    connection.set("scm:git:git://github.com/algorand/go-mobile-algorand-sdk.git")
                    developerConnection.set("scm:git:ssh://github.com/algorand/go-mobile-algorand-sdk.git")
                    url.set("https://github.com/algorand/go-mobile-algorand-sdk")
                }
            }
        }
    }

    repositories {
        mavenLocal()
    }
}